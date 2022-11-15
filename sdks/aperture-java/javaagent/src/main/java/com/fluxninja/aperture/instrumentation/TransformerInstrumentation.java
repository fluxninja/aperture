package com.fluxninja.aperture.instrumentation;

import net.bytebuddy.agent.builder.AgentBuilder;
import net.bytebuddy.description.type.TypeDescription;
import net.bytebuddy.matcher.ElementMatcher;

public interface TransformerInstrumentation {
    ElementMatcher<? super TypeDescription> getType();

    AgentBuilder.Transformer getTransformer();
}
