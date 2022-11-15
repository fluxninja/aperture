package com.fluxninja.aperture.instrumentation.armeria;

import com.fluxninja.aperture.instrumentation.InstrumentationModule;
import com.fluxninja.aperture.instrumentation.TransformerInstrumentation;

import java.util.List;

public class ArmeriaInstrumentationModule implements InstrumentationModule {
    @Override
    public List<TransformerInstrumentation> getTransformers() {
        return List.of(new ArmeriaClientInstrumentation(), new ArmeriaServerInstrumentation());
    }
}
