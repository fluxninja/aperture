package com.fluxninja.aperture.instrumentation;

import com.fluxninja.aperture.instrumentation.armeria.ArmeriaInstrumentationModule;

import com.fluxninja.aperture.instrumentation.netty.NettyInstrumentationModule;
import net.bytebuddy.agent.builder.AgentBuilder;

import java.lang.instrument.Instrumentation;
import java.util.List;

public class ApertureInstrumentationAgent {
    private static final List<InstrumentationModule> modules = List.of(
            new ArmeriaInstrumentationModule(),
            new NettyInstrumentationModule());

    public static void premain(String agentArgs, Instrumentation inst) {
        AgentBuilder agentBuilder = new AgentBuilder.Default()
                .with(new AgentBuilder.InitializationStrategy.SelfInjection.Eager());

        for (InstrumentationModule module : modules) {
            for (TransformerInstrumentation tf : module.getTransformers()) {
                agentBuilder = agentBuilder
                        .type(tf.getType())
                        .transform(tf.getTransformer());
            }
        }
        agentBuilder.installOn(inst);
    }
}
